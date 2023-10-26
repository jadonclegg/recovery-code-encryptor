import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EncryptionPageComponent } from './encryption-page.component';

describe('EncryptionPageComponent', () => {
  let component: EncryptionPageComponent;
  let fixture: ComponentFixture<EncryptionPageComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [EncryptionPageComponent]
    });
    fixture = TestBed.createComponent(EncryptionPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
